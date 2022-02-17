const MessageHandler = require("./MessageHandler/MessageHandler.js");
const fs = require("fs");

(async () => {
    const data = fs.readFileSync('machines.json','utf8');
    MessageHandler.sendMessage("New_Machines",data);
})();
  