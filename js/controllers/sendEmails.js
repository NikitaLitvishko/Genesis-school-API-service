const axios = require('axios');
const nodemailer = require('nodemailer');
const fs = require('fs');
const { EMAIL_HOST, EMAIL_HOST_PASSWORD } = require('../config');

module.exports = async function (req, res) {
  try {
    const current_rate = await axios.get('http://localhost:3000/api/rate');
    const transporter = nodemailer.createTransport({
      host: 'smtp.gmail.com',
      port: 465,
      secure: true,
      auth: {
        user: EMAIL_HOST,
        pass: EMAIL_HOST_PASSWORD,
      },
    });
    const emails = fs
      .readFileSync('subscribed_emails.txt', 'utf-8')
      .split(',');
    const messages = emails.map((email) => {
      return {
        from: `Gses2 <${EMAIL_HOST}>`,
        to: email,
        subject: 'Курс BTC/UAH',
        text: `Актуальний курс BTC до UAH: ${current_rate.data}`,
      };
    });

    for (const message of messages) {
      const info = await transporter.sendMail(message);
    }

    return res.send('Emails have been sent');
  }
  catch (err) {
    console.log(err);
  }
};
