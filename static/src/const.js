let host = window.location.host;
let ws = 'ws';
let http = 'http';
let siteUrl = (window.location !== window.parent.location) ? document.referrer : document.location.href;

if (NODE_ENV !== 'production') {
  host = 'localhost:8086';
} else {
  // ws = 'wss';
  // http = 'https';
}

// host = 'localhost:8086'; //localhost:8086, veliri.ru
// http = 'http';
// ws = 'ws';

export const urls = {
  authTokenKey: 'veliri-auth-token',
  socketURL: ws + "://" + host + "/socket",
  loginURL: http + "://" + host + "/api/login",
  siteUrl: siteUrl,
};

if (NODE_ENV === 'production') {
  console.log('There is Production mode');
} else {
  console.log('There is Development mode');
}

console.log('parent url: ', urls.siteUrl)

export default {urls}

