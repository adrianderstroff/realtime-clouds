#version 410 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec3 normal;

uniform mat4 M, V, P;

out Vertex {
    vec3 pos;
    vec3 normal;
} vertex;

void main(){
    gl_Position = P * V * M * vec4(pos, 1.0);
    vertex.pos = pos;
    vertex.normal = normal;
}