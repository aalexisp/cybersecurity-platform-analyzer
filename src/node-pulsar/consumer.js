const MessageHandler = require("./MessageHandler/MessageHandler.js");
const fs = require("fs");

(async () => {
    const message = await MessageHandler.receiveMessage("New_Machines","test");
    console.log(message.toString())
})();
  