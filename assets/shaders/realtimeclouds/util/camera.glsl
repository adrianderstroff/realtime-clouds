struct Camera {
    vec3  pos;      // position
    mat4  V;        // view matrix
    mat4  P;        // projection matrix
    float fov;      // vertical field of view
    float aspect;   // aspect ratio width/height
};