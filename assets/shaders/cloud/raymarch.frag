#version 430

//---------------------------------------------------------------------------------------//
// constants                                                                             //
//---------------------------------------------------------------------------------------//
const float PI           = 3.1415926535897932384626433832795;
const float PI_2         = 2*PI;
const float DEG_2_RAD    = PI/180.0;
const int   CLOUD_DOMAIN = 8000;
const float BASE_DENSITY = 0.4;
const float EXPOSURE     = 2.5;
const float INV_GAMMA    = 1.0 / 2.2;

//---------------------------------------------------------------------------------------//
// datastructures                                                                        //
//---------------------------------------------------------------------------------------//
struct Ray {
    vec3 o;
    vec3 dir;
};

struct Sphere {
    vec3  o;
    float r;
};

//---------------------------------------------------------------------------------------//
// textures                                                                              //
//---------------------------------------------------------------------------------------//
layout(binding = 0) uniform sampler3D cloudBaseTex;
layout(binding = 1) uniform sampler3D cloudDetailTex;
layout(binding = 2) uniform sampler2D turbulenceTex;
layout(binding = 3) uniform sampler2D cloudMapTex;

//---------------------------------------------------------------------------------------//
// uniforms                                                                              //
//---------------------------------------------------------------------------------------//
uniform vec3   cameraPos;
uniform mat4   M, V, P;
uniform float  width, height;
uniform float  fov = 45.0;
uniform vec3   sunPos = vec3(40000, -20000, 0);
uniform vec3   windDir = vec3(1, 0, 1);
uniform Sphere innerSphere = Sphere(vec3(0, 0, 0), 14000);
uniform Sphere outerSphere = Sphere(vec3(0, 0, 0), 40000);
uniform float  uTime = 0;
uniform float  uWindSpeed = 10;

//---------------------------------------------------------------------------------------//
// input                                                                                 //
//---------------------------------------------------------------------------------------//
in Vertex {
    vec2 uv;
} i;

//---------------------------------------------------------------------------------------//
// output                                                                                //
//---------------------------------------------------------------------------------------//
out vec4 fragColor;



//---------------------------------------------------------------------------------------//
// utility functions                                                                     //
//---------------------------------------------------------------------------------------//
// maps one value from one interval [inMin,inMax] to another interval [outMin, outMax]
float remap(in float val, in float inMin, in float inMax, in float outMin, in float outMax) {
    return (val - inMin)/(inMax - inMin) * (outMax - outMin) + outMin;
}
// clamps the input value to (inMin, inMax) and performs a remap
float clampRemap(in float val, in float inMin, in float inMax, in float outMin, in float outMax) {
    float clVal = clamp(val, inMin, inMax);
    return (clVal - inMin)/(inMax - inMin) * (outMax - outMin) + outMin;
}

// maps a position pos on a sphere with origin oSphere to uv-coordinates
// https://en.wikipedia.org/wiki/UV_mapping
vec2 sphereUV(in vec3 pos, in Sphere sphere) {
    vec3 dir = normalize(sphere.o - pos);
    float x = 0.5 + atan(dir.z, dir.x) / (2*PI);
    float y = 0.5 - asin(dir.y) / PI;
    return vec2(x, y);
}

vec2 boundUV(in vec3 pos, vec2 size) {
    float x = pos.x;
    float y = pos.z;
    
    while(x < 0.0) { x += size.x; }
    while(y < 0.0) { y += size.y; }

    x = mod(x, size.x) / size.x;
    y = mod(y, size.y) / size.y;

    return vec2(x, y);
}

// calculates the relative height of a sample position pos within the atmospher defined
// by the radii of the inner and outer sphere while the relative height is clamped
// between 0 and 1
float relativeHeight(in vec3 pos) {
    float len = length(pos - innerSphere.o);
    return clamp((len - innerSphere.r) / (outerSphere.r - innerSphere.r), 0.0, 1.0);
}

// swaps a and b
void swap(inout float a, inout float b) {
    float t = b;
    b = a;
    a = t;
}



//---------------------------------------------------------------------------------------//
// ray functions                                                                         //
//---------------------------------------------------------------------------------------//
// gets the ray origin and direction for the current fragment
Ray calcRay() {
    // calc image plane
    float ar = width / height;
    float angle = tan(fov/2 * DEG_2_RAD);
    vec2 imagePlane = (i.uv*2 - vec2(1)) * vec2(angle) * vec2(ar, 1);

    // extract camera space
    mat3 cameraToWorld = transpose(mat3(V));

    // ray origin is position of camera
    vec3 o = cameraPos;

    // transform direction with view matrix
    vec3 pLocal = vec3(imagePlane, -1);
    vec3 pWorld = cameraToWorld * pLocal;
    vec3 dir    = normalize(pWorld);

    return Ray(o, dir);
}

// Returns the intersection of the closer intersection point of a ray with a sphere
// @return a positive value if the sphere has been hit, while a negative value indicates no hit
// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-sphere-intersection
float intersectSphere(in Ray ray, in Sphere sphere) {
    vec3 L = ray.o - sphere.o;

    // calculate the determinant to check for intersection
    float a = dot(ray.dir, ray.dir);
    float b = 2.0 * dot(ray.dir, L);
    float c = dot(L, L) - (sphere.r * sphere.r);
    float d = b*b - 4*a*c;

    // determine the intersection parameters t0, t1
    float t0, t1;
    if (d < 0.0) {
        return -1.0;
    } else if (d == 0) {
        t0 = -0.5 * b/a;
        t1 = t0;
    } else {
        float q = (b > 0) ? 
            -0.5 * (b + sqrt(d)) :
            -0.5 * (b - sqrt(d));
        t0 = q / a;
        t1 = c / q;
    }
    // t0 should be the smaller value
    if(t0 > t1) swap(t0, t1);

    // check if t0 and or t1 are negative, if t0 is negative then use t1
    // instead and if both are negative then they are behind the ray origin
    if(t0 < 0) {
        t0 = t1;
        if(t0 < 0) return -1;
    }

    return t0;
}



//---------------------------------------------------------------------------------------//
// height gradient                                                                       //
//---------------------------------------------------------------------------------------//
float cloudRemap(float h, float a, float b, float c) {
    return clampRemap(h, 0, a, 0, 1) * clampRemap(h, b, c, 1, 0);
}
vec3 lerp3(float t) {
    float x = clamp(1 - t*2, 0, 1);
    float z = clamp((t-0.5)*2, 0, 1);
    float y = 1 - x - z;
    return vec3(x, y, z);
}
float heightGradient(float cloudType, float h) {
    // calc cloud gradients
    float a = cloudRemap(h, 0.1, 0.2, 0.3);
    float b = cloudRemap(h, 0.2, 0.3, 0.5);
    float c = cloudRemap(h, 0.1, 0.7, 0.8);
    // calc weights
    vec3 weights = lerp3(h);
    float cloudGradient = a*weights.x + b*weights.y + c*weights.z;
    return cloudGradient * h;
}



//---------------------------------------------------------------------------------------//
// cloud density                                                                         //
//---------------------------------------------------------------------------------------//
vec4 sampleCloudMap(in vec3 pos) {
    vec3 p = pos / CLOUD_DOMAIN;
    return texture(cloudMapTex, p.xz);
}

float sampleCloudBase(in vec3 pos, float h) {
    // grab value from texture
    vec3 p = pos / CLOUD_DOMAIN;
    vec4 cloudBase = texture(cloudBaseTex, vec3(p.xz, h));

    // calculate low freq fbm
    float lowFreqNoise = cloudBase.r;
    float highFreqNoise = dot(cloudBase.gba, vec3(0.625, 0.25, 0.125));
    float baseDensity = clampRemap(lowFreqNoise, highFreqNoise, 1.0, 0.0, 1.0);
    baseDensity = clamp(baseDensity, 0, 1);

    // extract cloud information
    vec4  cloudInfo = sampleCloudMap(pos);
    float coverage  = cloudInfo.r;
    float cloudType = cloudInfo.b;

    // apply height gradient
    baseDensity *= heightGradient(cloudType, h);
    baseDensity *= coverage;

    return baseDensity;
}

// returns the cloud density at the specified position
float density(in vec3 pos, in float height) {
    // offset position by wind direction and strength
    vec3 newPos = pos;
    
    // get base density
    float baseDensity = sampleCloudBase(pos, height) * BASE_DENSITY;

    // perform expensive sampling if ray is within cloud
    if(baseDensity > 0.0) {
        // calculate cloud detail density and erode it from
        // the base density
        // TODO
    }
    
    return baseDensity;
}



//---------------------------------------------------------------------------------------//
// cone sampling                                                                         //
//---------------------------------------------------------------------------------------//
// calculate cone sampling similar to Meteoros
// with the last sample being 3 times the 
// unit cone size
void coneSampling(inout vec3 coneSamples[5]) {
    coneSamples[0] = vec3(0.1, 0.1, 0.0);
    coneSamples[1] = vec3(0.2, 0.0, 0.2);
    coneSamples[2] = vec3(0.4,-0.4, 0.0);
    coneSamples[3] = vec3(0.8, 0.0,-0.8);
    coneSamples[4] = vec3(3.0, 0.0, 0.0);
}

// creates a rotation matrix that rotates vector a onto 
// vector b by performing a 2D rotation on the plane
// with normal a x b
// source: https://math.stackexchange.com/questions/180418/calculate-rotation-matrix-to-align-vector-a-to-vector-b-in-3d
mat3 coneRotationMatrix(in vec3 coneDir, in vec3 lightDir) {
    vec3 a = normalize(coneDir);
    vec3 b = normalize(lightDir);
    return mat3(
        dot(a, b)          , -length(cross(a, b)), 0.0,
        length(cross(a, b)),  dot(a, b)          , 0.0,
        0.0                ,  0.0                , 1.0
    );
}

float sampleConeDensity(in vec3 pos, in vec3 coneSamples[5]) {
    // calculate orientation towards the sun
    vec3 toLight = normalize(sunPos - pos);
    mat3 coneRot = coneRotationMatrix(vec3(1, 0, 0), toLight);

    // sample density along cone
    float cloudDensity= 0.0;
    for(int i = 0; i < 5; i++) {
        vec3 coneSamplePos = pos + (coneRot * coneSamples[i] * 10);
        float coneSampleHeight = relativeHeight(coneSamplePos);
        cloudDensity += sampleCloudBase(coneSamplePos, coneSampleHeight);
    }
    return cloudDensity;
}



//---------------------------------------------------------------------------------------//
// lighting                                                                              //
//---------------------------------------------------------------------------------------//
float beerLambert(float d) {
    return max(exp(-d), 0.7*exp(-0.25*d));
}

// the Henyey-Greenstein formula is used to approximate Mie scattering which is too
// computationally expensive
float henyeyGreenstein(float cosAngle, float eccentricity) {
    float eccentricity2 = eccentricity*eccentricity;
    return ((1.0 - eccentricity2) / pow(1.0 + eccentricity2 - 2.0*eccentricity*cosAngle, 3.0/2.0)) / (4*PI);
}

float inscatter(float d, float height) {
    // attenuation along the in-scatter path
    float depthProbability = 0.05 + pow(d, clampRemap(height, 0.3, 0.85, 0.5, 2.0));
    // relax the attenuation over hight
    float verticalProbability = pow(clampRemap(height, 0.07, 0.14, 0.1, 1.0), 0.8);
    // both of those effects model the in-scatter probability
    float inScatterProbability = depthProbability * verticalProbability;
    return inScatterProbability;
}

float radiance(vec3 pos, in vec3 coneSamples[5], float height, float silverIntensity, float silverSpread) {
    // get depth from cone sampling
    float d = sampleConeDensity(pos, coneSamples);
    
    // we wanna calculate wether we look towards the sun or
    // away from it and depending on the adapt the anisotropic
    // scattering defined by the henyey-greenstein function
    vec3 toSun = sunPos - pos;
    vec3 toEye   = cameraPos - pos;
    float cosAngle = dot(normalize(toSun), normalize(toEye));

    // two henyey-greenstein functions were combined to have
    // a highlight around the sun but also to retain silver
    // lining highlights on the clouds that are 90 degrees
    // away from the sun
    float eccentricity = 0.6;
    float HG1 = henyeyGreenstein(cosAngle, eccentricity);
    float HG2 = henyeyGreenstein(cosAngle, 0.99 - silverSpread);
    float hg = max(HG1, silverIntensity*HG2);

    // calculate the out-scattering and obsorption of light
    // that travels through the cloud
    // this effect is direction dependent, it gets stronger
    // the further the camera looks away from the sun
    // TODO: add view dependent scaling
    float bl = beerLambert(d);

    // calcluate the in-scattering contribution which
    // creates dark edges since less in-scattering is
    // taking place or vice versa a lot of in-scattering
    // happens in thicker parts of the cloud
    // in addition are the bottoms of clouds darker since
    // below them is no medium that scatters the light 
    // back into the cloud bottom
    float is = inscatter(d, height);

    // the radiance at the position pos is the combination
    // of attenuation, absorption/out-scattering and 
    // in-scattering
    return bl * hg * is;
}



//---------------------------------------------------------------------------------------//
// tone mapping                                                                          //
//---------------------------------------------------------------------------------------//
// Uncharted 2 tone mapping
const float A = 0.15;
const float B = 0.50;
const float C = 0.10;
const float D = 0.20;
const float E = 0.02;
const float F = 0.30;
vec3 unchartedToneMapping(in vec3 x) {
    return ((x*(A*x + C*B) + D*E) / (x*(A*x + B) + D*F)) - E/F;
}
vec3 toneMapping(in vec3 x) {
    vec3 color = unchartedToneMapping(x*EXPOSURE);
    vec3 white = vec3(100);
    color *= 1.0 / unchartedToneMapping(white);
    return pow(color, vec3(INV_GAMMA));
}



//---------------------------------------------------------------------------------------//
// entry point                                                                           //
//---------------------------------------------------------------------------------------//
void main() {
    // determine ray direction
    Ray ray = calcRay();

    // get start and end points of the horizon
    float tInner = intersectSphere(ray, innerSphere);
    float tOuter = intersectSphere(ray, outerSphere);

    // step size
    float stepSize = (tOuter - tInner) / 30.0;

    // setup cone sample
    vec3 coneSamples[5];
    coneSampling(coneSamples);

    // setup ray marching variables
    float light         = 0.0;
    float alpha         = 0.0;
    float transmittance = 1.0;
    vec3  accumColor    = vec3(0.0);

    // perform ray marching
    float t = tInner;
    while(t <= tOuter) {
        // get position within cloud layer
        vec3 pos = ray.o + ray.dir*t;
        float h = remap(t, tInner, tOuter, 0, 1);

        // offset by time
        pos += uWindSpeed *uTime * h;

        // calculate density
        float d = density(pos, h);
        alpha += d;

        // do light calculations
        if(d > 0) {
            float lightEnergy = radiance(pos, coneSamples, h, 0.7, 1.0);
            transmittance = mix(transmittance, lightEnergy, (1.0 - alpha));
            accumColor += vec3(transmittance);
        }

        // advance ray position
        t += stepSize;

        // early ray termination
        if(alpha > 1.0) {
            alpha = 1.0;
            break;
        }
    }

    // Blend and fade out clouds into the horizon (CHANGE THIRD PARAM IN REMAP)
    const float cloudFadeOutPoint = 0.06f;
    alpha *= smoothstep(0.0, 1.0, min(1.0, remap(ray.dir.y, cloudFadeOutPoint, 0.2f, 0.0f, 1.0f)));

    // calculate light color from light
    vec3 atmosphereColor = vec3(0, 0, 1);
    fragColor = vec4(mix(atmosphereColor, accumColor, alpha), 1.0);

    // apply tone mapping
    fragColor.xyz = toneMapping(fragColor.xyz);
}