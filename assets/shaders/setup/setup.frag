#version 410 core

in Vertex {
    vec3 pos;
} i;

out vec3 color;

void main() {
    color = i.pos;
}