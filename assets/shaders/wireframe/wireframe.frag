#version 410

in Vertex {
    vec3 barycoord;
} i;

uniform float width = 1.5;

layout(location = 0) out vec4 fragColor;

float edgeFactor(){
    vec3 d = fwidth(i.barycoord);
    vec3 a3 = smoothstep(vec3(0.0), d*width, i.barycoord);
    return min(min(a3.x, a3.y), a3.z);
}

void main() {
    fragColor = vec4(0.0, 0.0, 0.0, (1.0-edgeFactor())*0.95);
}