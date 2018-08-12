#version 410 core

in Vertex {
    vec2 uv;
} i;

uniform sampler2D tex;

out vec3 color;

void main(){
    color = texture(tex, i.uv).xyz;
}