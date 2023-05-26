const express = require('express');
const app = express();
const port = 3000;
const {
  getBitcoinExchangeRate,
  subscribeEmail,
  sendEmails,
} = require('./controllers/index');
const { json } = require('express');

app.use(json());
app.use(express.urlencoded({ extended: true }));

app.get('/api/rate', getBitcoinExchangeRate);
app.get('/api/subscribe', subscribeEmail);
app.get('/api/sendEmails', sendEmails);

app.listen(port, () => {
  console.log(`Server is listening on port ${port}`);
});
