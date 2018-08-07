#version 430

in VertexOut {
    vec3  position;
    vec2  uv;
    vec3  normal;
    float texID;
} i;

// textures
layout(binding = 0) uniform sampler2D grassAlpha;
// write vec3 color
layout(location = 0) out vec3 fragColor;

void main() {
    // discard transparent pixels
    if(texture(grassAlpha, i.uv).x < 0.5 && i.texID != 0) {
        discard;
    }

    // select color depending on LOD
    fragColor = vec3(0, 0, 1);
    if      (i.texID == 3) { fragColor = vec3(1, 0, 0); } 
    else if (i.texID == 2) { fragColor = vec3(1, 1, 0); } 
    else if (i.texID == 1) { fragColor = vec3(0, 1, 0); }

    // ambient occlusion
    if (i.texID > 0) {
        fragColor = mix(vec3(0, 0, 0), fragColor, 1-i.uv.y);
    }
}