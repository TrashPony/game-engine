//  NODE_ENV определяется при сборке в файле package.json, раздел scripts {}

let host = window.location.host;
let ws = 'ws';
let http = 'http';

if (NODE_ENV !== 'production') {
  host = 'localhost:8080';
} else {
  ws = 'wss';
  http = 'https';
}

export const urls = {
  socketURL: ws + "://" + host + "/socket",
  regURL: http + "://" + host + "/api/registration",
  loginURL: http + "://" + host + "/api/login",
  vkOAuth: http + "://" + host + "/api/vk-oauth",
  vkOAuthUrl: http + "://" + host + "/api/vk-get-oauth-url",
  vkAppLogin: http + "://" + host + "/api/vk-app-login",
};

if (NODE_ENV === 'production') {
  console.log('There is Production mode');
} else {
  console.log('There is Development mode');
}

export default {urls}

