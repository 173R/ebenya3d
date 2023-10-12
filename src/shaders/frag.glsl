#version 330
in vec3 vertexColor;
out vec4 fragColour;

uniform vec4 ourColor;

void main() {
    fragColour = vec4(vertexColor, 1.0);
}