const axios = require("axios");
const hmacsha1 = require('hmacsha1');
const base64 = value => Buffer.from(value).toString("base64")
function makeid(length) {
   var result           = '';
   var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
   var charactersLength = characters.length;
   for ( var i = 0; i < length; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
   return result;
}
const en = str => encodeURIComponent(str);

const { CONSUMER_KEY, CONSUMER_SECRET, APP_TOKEN, APP_SECRET, TWEET_MSG: msg } = process.env;
const nonce = base64(makeid(32));
const URL = "https://api.twitter.com/1.1/statuses/update.json?status=" + en(msg)

const SIGING_KEY = `${en(CONSUMER_SECRET)}&${en(APP_SECRET)}`;

const headers = [
  {name: "oauth_consumer_key", value: CONSUMER_KEY},
  {name: "oauth_nonce", value: nonce},
  {name: "oauth_signature_method", value: "HMAC-SHA1"},
  {name: "oauth_timestamp", value: parseInt(new Date().getTime()/1000, 10)},
  {name: "oauth_version", value: "1.0"},
  {name: "oauth_token", value: APP_TOKEN},
  {name: "status", value: msg}
];

const base_str = `POST&${en("https://api.twitter.com/1.1/statuses/update.json")}&${en(headers.map(e => ({name: en(e.name), value: en(e.value)})).sort((a,b) => {
  if(a.name < b.name) return -1;
  if(a.name > b.name) return 1;
  return 0;
}).map(entry => `${entry.name}=${entry.value}`).join("&"))}`

const signature = hmacsha1(SIGING_KEY, base_str);
headers.push({name: "oauth_signature", value: signature})
const headersObj = {
  Authorization: `OAuth ${headers.filter(e => e.name !== "status").map(e => `${en(e.name)}="${en(e.value)}"`).sort((a,b) => {
  if(a < b) return -1;
  if(a > b) return 1;
  return 0;
}).join(", ")}`
}
axios.post(URL, null, {
  headers: headersObj
}).then(res => {
  console.log(res.data)

}).catch(err => {
  console.log(err.response.data)
})
