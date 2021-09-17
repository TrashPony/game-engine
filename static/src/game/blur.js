import Phaser from "phaser";

const Blur = new Phaser.Class({

  Extends: Phaser.Renderer.WebGL.Pipelines.SinglePipeline,

  initialize:

    function CustomPipeline(game) {
      Phaser.Renderer.WebGL.Pipelines.SinglePipeline.call(this, {
        game: game,
        fragShader: `
            precision mediump float;
            //in attributes from our vertex shader
            varying vec4 outColor;
            varying vec2 outTexCoord;
            //declare uniforms
            uniform sampler2D u_texture;
            uniform vec2 dir;
            uniform float alpha;

            float resolution = 800.0;
            float radius = 3.0;

            void main() {
                //this will be our RGBA sum
                vec4 sum = vec4(0.0);
                //our original texcoord for this fragment
                vec2 tc = outTexCoord;
                //the amount to blur, i.e. how far off center to sample from
                //1.0 -> blur by one pixel
                //2.0 -> blur by two pixels, etc.
                float blur = radius/resolution;
                //the direction of our blur
                //(1.0, 0.0) -> x-axis blur
                //(0.0, 1.0) -> y-axis blur
                float hstep = 2.0;
                float vstep = 2.0;

                //apply blurring, using a 9-tap filter with predefined gaussian weights,
                sum += texture2D(u_texture, vec2(tc.x - 4.0*blur*hstep, tc.y - 4.0*blur*vstep)) * 0.0162162162;
                sum += texture2D(u_texture, vec2(tc.x - 3.0*blur*hstep, tc.y - 3.0*blur*vstep)) * 0.0540540541;
                sum += texture2D(u_texture, vec2(tc.x - 2.0*blur*hstep, tc.y - 2.0*blur*vstep)) * 0.1216216216;
                sum += texture2D(u_texture, vec2(tc.x - 1.0*blur*hstep, tc.y - 1.0*blur*vstep)) * 0.1945945946;
                sum += texture2D(u_texture, vec2(tc.x, tc.y)) * 0.2270270270;
                sum += texture2D(u_texture, vec2(tc.x + 1.0*blur*hstep, tc.y + 1.0*blur*vstep)) * 0.1945945946;
                sum += texture2D(u_texture, vec2(tc.x + 2.0*blur*hstep, tc.y + 2.0*blur*vstep)) * 0.1216216216;
                sum += texture2D(u_texture, vec2(tc.x + 3.0*blur*hstep, tc.y + 3.0*blur*vstep)) * 0.0540540541;
                sum += texture2D(u_texture, vec2(tc.x + 4.0*blur*hstep, tc.y + 4.0*blur*vstep)) * 0.0162162162;

                sum.x = sum.x*0.67;
                sum.y = sum.y*0.35;
                sum.z = sum.z*0.15;

                gl_FragColor = sum*alpha;
            }
       `,
        uniforms: [
          'alpha'
        ]
      });
    }
});

'rgb(181,104,65)'

export {Blur}
