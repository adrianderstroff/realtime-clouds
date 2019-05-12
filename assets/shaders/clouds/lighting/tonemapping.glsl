// camera parameters
const float EXPOSURE     = 2.5;
const float INV_GAMMA    = 1.0 / 2.2;

// Uncharted 2 tone mapping parameters
const float A = 0.15;
const float B = 0.50;
const float C = 0.10;
const float D = 0.20;
const float E = 0.02;
const float F = 0.30;

// Uncharted 2 tone mapping function
vec3 unchartedToneMapping(in vec3 x) {
    return ((x*(A*x + C*B) + D*E) / (x*(A*x + B) + D*F)) - E/F;
}

// General tone mapping
vec3 toneMapping(in vec3 x) {
    vec3 color = unchartedToneMapping(x*EXPOSURE);
    vec3 white = vec3(100);
    color *= 1.0 / unchartedToneMapping(white);
    return pow(color, vec3(INV_GAMMA));
}