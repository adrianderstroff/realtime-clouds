#version 400

in vec3 tex;
uniform samplerCube cubeMap;
out vec3 color;

void main(){
    color = vec3(texture(cubeMap, tex));
}