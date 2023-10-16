#version 330
in vec3 vertexColor;
in vec2 TexCoord;

out vec4 fragColor;

//uniform vec4 ourColor;

uniform sampler2D ourTexture;

void main() {
    //fragColour = vec4(vertexColor, 1.0);
    fragColor = texture(ourTexture, TexCoord);
}

