new Vue({
    el: '#app',
    data:{
        ws:null, // websocket
        newMsg:'', // holds new message to be sent
        chatContent: '', // running list of chat messages displayed on the screen
        email: null, // email address of user
        username: null, // user name
        joined: false // true if email and username have  been filled
    },
    created: function(){
        var self=this;
        this.ws = new WebSocket('ws://'+ window.location.host+'/ws');
        this.ws.addEventListener('message', function(e){
            var msg= JSON.parse(e.data);
            self.chatContent+= '<div class="chip">'+
            '<img src="'+ self.gravatarURL(msg.email) +'">'+
            msg.username + '</div>' +
            emojione.toImage(msg.message)+
            '<br/>'; 

            var element=document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight;// auto scroll to bottom


        });
    },

    methods:{
        send: function(){
            if(this.newMsg!= ''){
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        // strip html data
                        message: $('<p>').html(this.newMsg).text()
                    })
                );
                this.newMsg=''; // reset newMsg
            }//endif
        },
        join:function(){
            if(!this.email){
                Materialize.toast('You must enter an email', 2000);
                return;
            }
            if(!this.username){
                Materialize.toast('You must choose an username', 2000);
                return;
            }

            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined=true;
        },
        gravatarURL: function(email){
            return 'http://www.gravatar.com/avatar/'+ CryptoJS.MD5(email);
        }
    }
});