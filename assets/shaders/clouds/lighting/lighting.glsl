// returns the cloud depth around a position
float getLODCloudDepth(in vec3 pos, float ds) {
    // get sample positions around the position
    vec3 p1 = pos + vec3(ds, 0, 0);
    vec3 p2 = pos - vec3(ds, 0, 0);
    vec3 p3 = pos + vec3( 0,ds, 0);
    vec3 p4 = pos - vec3( 0,ds, 0);
    vec3 p5 = pos + vec3( 0, 0,ds);
    vec3 p6 = pos - vec3( 0, 0,ds);

    // get relative height
    float h1 = relativeHeight(p1);
    float h2 = relativeHeight(p2);
    float h3 = relativeHeight(p3);
    float h4 = relativeHeight(p4);
    float h5 = relativeHeight(p5);
    float h6 = relativeHeight(p6);

    // sum depth for all sampled positions
    float sum = 0;
    sum += density(p1, h1);
    sum += density(p2, h2);
    sum += density(p3, h3);
    sum += density(p4, h4);
    sum += density(p5, h5);
    sum += density(p6, h6);
    sum /= 6;

    return sum;
}

float densityToLight(in vec3 pos, in int samples) {
    vec3 toLight = normalize(uSunPos - pos);
    float stepSize = (outerHeight - innerHeight) / samples;

    float accum = 0;
    for(int i = 0; i < samples; i++) {
        vec3 p = pos + toLight * (stepSize * i);
        float h = relativeHeight(p);
        accum += density(p, h);
    }

    return accum / samples;
}

// models the transmittance over depth of the light ray
float beerLambert(float d) {
    return max(exp(-d), 0.7*exp(-0.25*d));
}

// the Henyey-Greenstein formula is used to approximate Mie scattering which is too
// computationally expensive
float henyeyGreenstein(float cosAngle, float eccentricity) {
    float eccentricity2 = eccentricity*eccentricity;
    return ((1.0 - eccentricity2) / pow(1.0 + eccentricity2 - 2.0*eccentricity*cosAngle, 3.0/2.0)) / (4*PI);
}

// in scatter describes the light contribution from all directions that scatter
// towards the eye ray. this results in brighter parts where the cloud is thicker
// as there are higher in-scatter probabilities. in addition does the cloud gets
// less contribution at the bottom of a cloud as there is nothing below the cloud
// that can scatter the light back into the cloud 
float inscatter(float d, float height) {
    // attenuation along the in-scatter path
    float depthProbability = 0.05 + pow(d, clampRemap(height, 0.3, 0.85, 0.5, 2.0));
    // relax the attenuation over height
    float verticalProbability = pow(clampRemap(height, 0.07, 0.14, 0.1, 1.0), 0.8);
    // both of those effects model the in-scatter probability
    float inScatterProbability = depthProbability * verticalProbability;
    return inScatterProbability;
}

// calculates the light energy for the current sample position consisting of
// in-scatter, out-scattering and absorption
float radiance(vec3 pos, in vec3 coneSamples[5], float height, float silverIntensity, float silverSpread) {    
    // we wanna calculate wether we look towards the sun or
    // away from it and depending on the adapt the anisotropic
    // scattering defined by the henyey-greenstein function
    vec3 toSun = uSunPos - pos;
    vec3 toEye   = uCameraPos - pos;
    float cosAngle = dot(normalize(toSun), normalize(toEye));

    // two henyey-greenstein functions were combined to have
    // a highlight around the sun but also to retain silver
    // lining highlights on the clouds that are 90 degrees
    // away from the sun
    float eccentricity = 0.6;
    float HG1 = henyeyGreenstein(cosAngle, eccentricity);
    float HG2 = henyeyGreenstein(cosAngle, 0.99 - silverSpread);
    float hg = max(HG1, silverIntensity*HG2);

    // calculate the out-scattering and absorption of light
    // that travels from the sun through the cloud to the 
    // current position. this is simplified by using the 
    // cone sampled depth instead of calculating a light ray
    // this effect is direction dependent, it gets stronger
    // the further the camera looks away from the sun
    // TODO: add view dependent scaling
    float d = sampleConeDensity(pos, coneSamples);
    float bl = beerLambert(d);

    // calcluate the in-scattering contribution which
    // creates dark edges since less in-scattering is
    // taking place or vice versa a lot of in-scattering
    // happens in thicker parts of the cloud
    // in addition are the bottoms of clouds darker since
    // below them is no medium that scatters the light 
    // back into the cloud bottom
    float lodd = getLODCloudDepth(pos, 1);
    float is = inscatter(lodd, height);

    // the radiance at the position pos is the combination
    // of attenuation, absorption/out-scattering and 
    // in-scattering
    return bl * hg * is;
}

vec3 radianceSimple(vec3 pos, float height, float silverIntensity, float silverSpread) {    
    // we wanna calculate wether we look towards the sun or
    // away from it and depending on the adapt the anisotropic
    // scattering defined by the henyey-greenstein function
    vec3 toSun = normalize(uSunPos - pos);
    vec3 toEye = normalize(uCameraPos - pos);
    float cosAngle = dot(toSun, toEye);

    // two henyey-greenstein functions were combined to have
    // a highlight around the sun but also to retain silver
    // lining highlights on the clouds that are 90 degrees
    // away from the sun
    float eccentricity = 0.6;
    float HG1 = henyeyGreenstein(cosAngle, eccentricity);
    float HG2 = henyeyGreenstein(cosAngle, 0.99 - silverSpread);
    float hg = max(HG1, silverIntensity*HG2);

    // calculate the out-scattering and absorption of light
    // that travels from the sun through the cloud to the 
    // current position. this is simplified by using the 
    // cone sampled depth instead of calculating a light ray
    // this effect is direction dependent, it gets stronger
    // the further the camera looks away from the sun
    // TODO: add view dependent scaling
    float d = densityToLight(pos, 10);
    float bl = beerLambert(d);

    // calcluate the in-scattering contribution which
    // creates dark edges since less in-scattering is
    // taking place or vice versa a lot of in-scattering
    // happens in thicker parts of the cloud
    // in addition are the bottoms of clouds darker since
    // below them is no medium that scatters the light 
    // back into the cloud bottom
    float lodd = getLODCloudDepth(pos, 1);
    float is = inscatter(lodd, height);

    // the radiance at the position pos is the combination
    // of attenuation, absorption/out-scattering and 
    // in-scattering
    vec3 sunLight = uSunColor * (bl * hg);
    vec3 ambientLight = uAmbientColor * is;
    return mix(sunLight, ambientLight, height);
}