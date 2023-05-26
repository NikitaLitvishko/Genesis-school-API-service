const axios = require('axios');
const { API_KEY, BTC_ID, UAH_ID } = require('../config');


module.exports = async function (req, res) {
  try {
    const response = await axios.get(
      'https://pro-api.coinmarketcap.com/v2/tools/price-conversion',
      {
        params: {
          amount: 1,
          id: BTC_ID,
          convert_id: UAH_ID,
        },
        headers: {
          'X-CMC_PRO_API_KEY': API_KEY,
        },
      }
    );
    return res.json(response.data.data.quote['2824'].price);
  } catch (err) {
    return res.status(400).send(err);
  }
}