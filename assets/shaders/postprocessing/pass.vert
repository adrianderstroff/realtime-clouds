#version 410

layout(location = 0) in vec3 vert;
layout(location = 1) in vec2 uv;

out VertexOut {
    vec2 uv;
}o;

void main() {
    gl_Position = vec4(vert - vec3(0, 0, 1), 1.0);
    o.uv = uv;
}