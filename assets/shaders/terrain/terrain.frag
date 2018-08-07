#version 430

in VertexOut {
    vec3 position;
    vec3 normal;
} i;

// phong shader uniforms
uniform vec3  cameraPos;
uniform vec3  lightDir;
uniform vec3  lightColor;
uniform float ambientIntensity;
uniform float diffuseIntensity;
// saturation distances
uniform float d1;
uniform float d2;

layout(location = 0) out vec3 fragColor;

void main() {
    fragColor = vec3(0.0, 0.0, 0.0);
}