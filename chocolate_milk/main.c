#include <windows.h>
#include <stdio.h>
#include <stdint.h>
#include <stddef.h>
#include <assert.h>

#define PLAYER_POINTER_OFFSET 0xCD87A8
#define PLAYER_GUID_OFFSET 0xCA1238

#define CURRENT_CLIENT_CONNECTION_OFFSET 0xC79CE0
#define OBJ_MANAGER_OFFSET 0x2ED0
#define FIRST_OBJ_OFFSET 0xAC
#define NEXT_OBJ_OFFSET 0x3C
#define OBJ_TYPE_OFFSET 0x14

typedef uint64_t u64;
typedef uint32_t u32;

#define MAX(x, y) (((x) > (y)) ? (x) : (y))

// --- Stretchy buffers ---
// Just a nice way to implment dynamic arrays. Instead of having something
// like a access function we can index the dynamic array with the same notation
// we index static ones.
typedef struct {
    size_t len;
    size_t cap;
    char buf[0];
} Buf_Hdr;

#define buf__hdr(b) ((b) ? (Buf_Hdr *)(((char *)b) - offsetof(Buf_Hdr, buf)) : NULL)
#define buf__fit(b, n) (((buf_len(b) + n) < buf_cap(b)) ? 0 : ((b) = buf__resize((b), buf_len(b)+(n), sizeof(*b))))

#define buf_len(b) ((b) ? buf__hdr(b)->len : 0)
#define buf_cap(b) ((b) ? buf__hdr(b)->cap : 0)

#define buf_push(b, e) (buf__fit((b), 1), (b)[buf__hdr(b)->len++] = (e))
#define buf_free(b) ((b) ? free(buf__hdr(b)), (b) = NULL : 0)

inline u64 read_u64(u64 addr) { return *(u64*)addr; }

inline u32 read_u32(u32 addr) { return *(u32*)addr; }

void die(const char* str) {
    fprintf(stderr, "[!] %s\n", str);
    exit(1);
}

void* xrealloc(void* ptr, size_t new_size) {
    ptr = realloc(ptr, new_size);
    if (!ptr) die("xrealloc failed");
    return ptr;
}

void* xmalloc(size_t size) {
    void* ptr = malloc(size);
    if (!ptr) die("xmalloc failed");
    return ptr;
}

void* buf__resize(void* ptr, size_t min_new_cap, size_t elem_size) {
    size_t new_cap = MAX(min_new_cap, 1 + buf_cap(ptr) * 2);
    size_t new_size = offsetof(Buf_Hdr, buf) + (new_cap * elem_size);

    Buf_Hdr* new_buf;
    if (ptr) {
        new_buf = xrealloc(buf__hdr(ptr), new_size);
    }
    else {
        new_buf = xmalloc(new_size);
        new_buf->len = 0;
    }
    new_buf->cap = new_cap;

    return new_buf->buf;
}

void buf_test() {
    int* arr = NULL;
    assert(buf_len(arr) == 0);
    assert(buf_cap(arr) == 0);

    for (int i = 0; i < 1024; i++) {
        buf_push(arr, i);
    }
    assert(buf_len(arr) == 1024);

    for (int i = 0; i < 1024; i++) {
        assert(i == arr[i]);
    }

    buf_free(arr);
    assert(buf_len(arr) == 0);
    assert(buf_cap(arr) == 0);
    assert(arr == NULL);
}

u64 get_local_player_guid() {
    return read_u64(PLAYER_GUID_OFFSET);
}

Vector3 get_object_pos() {
}

typedef enum {
    None,
    Item,
    Container,
    Unit,
    Player,
    GameObject,
    DynamicObject,
    Corpse
} Object_Type;

typedef struct { float x, y, z; } Vector3;

typedef struct {
    Object_Type type;
    u32 base_addr;
} Wow_Object;

Wow_Object* get_visible_objects() {
    Wow_Object* objects = NULL;

    u32 curr_obj_mgr = read_u32(CURRENT_CLIENT_CONNECTION_OFFSET);
    curr_obj_mgr = read_u32(curr_obj_mgr + OBJ_MANAGER_OFFSET);

    u32 curr_obj = read_u32(curr_obj_mgr + FIRST_OBJ_OFFSET);
    u32 obj_type = read_u32(curr_obj + OBJ_TYPE_OFFSET);
    printf("%d\n", obj_type);

    while (obj_type <= 7 && obj_type > 0) {
        buf_push(objects, ((Wow_Object){obj_type, curr_obj}));

        curr_obj = read_u32(curr_obj + NEXT_OBJ_OFFSET);
        obj_type = read_u32(curr_obj + OBJ_TYPE_OFFSET);
    }

    return objects;
}

void main_routine(HINSTANCE inst) {
    buf_test();

    // Alloc console
    AllocConsole();
    FILE* f;
    freopen_s(&f, "CONOUT$", "w", stdout);

    printf("starting up...");

    Wow_Object* objs = get_visible_objects();
    for (int i = 0; i < buf_len(objs); i++) {
        printf("BAsE: %x\tTYPE: %d\n", objs[i].base_addr, objs[i].type);
    }
    buf_free(objs);

    for (;;) {
        if (GetAsyncKeyState(VK_END) & 1) {
            break;
        }
        Sleep(100);
    }

    // Dealloc console
    fclose(f);
    FreeConsole();
    FreeLibraryAndExitThread(inst, 0);
}

int DllMain(HINSTANCE inst, DWORD reason, LPVOID reserved) {
    if (reason == DLL_PROCESS_ATTACH) {
        CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)main_routine, inst, 0, NULL);
    }

    return TRUE;
}
