const fs = require('fs');

module.exports = function (req, res) {
  console.log(req.body);
  let email = '';
  if (req.body) email = req.body.email;

  if (!email) {
    return res.status(400).send('Email is required');
  }

  try {
    const subscribedEmails = fs.readFileSync('subscribed_emails.txt', 'utf-8');

    if (subscribedEmails.includes(email)) {
      return res
        .status(409)
        .send(`The ${email} exists in subscriptions`);
    }

    const updatedSubscribedEmails = subscribedEmails ? `${subscribedEmails},${email}` : email;
    fs.writeFileSync('subscribed_emails.txt', updatedSubscribedEmails, 'utf-8');
    return res.send('E-mail has been added');
  } catch (err) {
    return res.status(500).send(`Error handling email: ${err}`);
  }
};
