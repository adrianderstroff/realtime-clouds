#version 410 core

#include "phong.glsl"

in Vertex {
    vec3 pos;
    vec3 normal;
} vertex;

uniform vec3 cameraPos;

out vec3 fragColor;

const vec3 lightPos = vec3(5, 5, -5);
const PhongMaterial mat = PhongMaterial(
    vec3(0.2, 0.2, 0.2), 0.2, 
    vec3(0.3, 0.3, 0.5), 0.3, 
    vec3(1, 1, 1), 0.2, 50
);

void main(){
    fragColor = vec3(0.3, 0, 0) + phong(vertex.pos, vertex.normal, cameraPos, lightPos, mat);
}