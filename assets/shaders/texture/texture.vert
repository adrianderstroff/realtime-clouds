#version 410 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 uv;

uniform mat4 M, V, P;

out Vertex {
    vec2 uv;
} o;

void main(){
    gl_Position = P * V * M * vec4(pos, 1.0);
    o.uv = uv;
}