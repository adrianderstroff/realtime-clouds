#version 410 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 uv;

out Vertex {
    vec2 uv;
} o;

void main() {
    gl_Position = vec4(pos.xz, 0 , 1);
    o.uv = uv;
}