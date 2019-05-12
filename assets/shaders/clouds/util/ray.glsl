struct Ray {
    vec3 o;
    vec3 dir;
};

// gets the ray origin and direction for the current fragment
Ray calcRay() {
    // calc image plane
    float ar = width / height;
    float angle = tan(radians(fov/2));
    vec2 imagePlane = (i.uv*2 - vec2(1)) * vec2(angle) * vec2(ar, 1);

    // extract camera space
    mat3 cameraToWorld = transpose(mat3(V));

    // ray origin is position of camera
    vec3 o = uCameraPos;

    // transform direction with view matrix
    vec3 pLocal = vec3(imagePlane, -1);
    vec3 pWorld = cameraToWorld * pLocal;
    vec3 dir    = normalize(pWorld);

    return Ray(o, dir);
}