// Vertex shader for protein atom rendering
// Instanced rendering: one draw call for all atoms

attribute vec3 aPosition;        // Sphere vertex position
attribute vec3 aInstancePosition; // Per-atom center
attribute vec4 aInstanceColor;    // Per-atom color
attribute float aInstanceRadius;  // Per-atom radius

uniform mat4 uProjectionMatrix;
uniform mat4 uViewMatrix;

varying vec4 vColor;
varying vec3 vNormal;

void main() {
    // Scale sphere by radius and translate to atom position
    vec3 worldPos = aPosition * aInstanceRadius + aInstancePosition;
    gl_Position = uProjectionMatrix * uViewMatrix * vec4(worldPos, 1.0);

    // Normal for lighting (sphere surface normal)
    vNormal = normalize(aPosition);
    vColor = aInstanceColor;
}
