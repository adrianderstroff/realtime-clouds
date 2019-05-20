#include "camera.glsl"

struct Ray {
    vec3 o;
    vec3 dir;
};

// gets the ray origin and direction for the current fragment
Ray calcRay(in vec2 uv, in Camera cam) {
    // calc image plane
    float angle = tan(radians(cam.fov/2));
    vec2 imagePlane = (uv*2 - vec2(1)) * vec2(angle) * vec2(cam.aspect, 1);

    // extract camera space
    mat3 cameraToWorld = transpose(mat3(cam.V));

    // ray origin is position of camera
    vec3 o = cam.pos;

    // transform direction with view matrix
    vec3 pLocal = vec3(imagePlane, -1);
    vec3 pWorld = cameraToWorld * pLocal;
    vec3 dir    = normalize(pWorld);

    return Ray(o, dir);
}