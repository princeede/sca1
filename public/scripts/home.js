(function (){
    var myId;
    buttonTestArr = $(".supportButton").click(function(event){
        myId = event.currentTarget.dataset["some"]
    });

    // var myForm =  document.querySelector("form")
    document.forms["supportForm"].addEventListener("submit", e => {
        var x = document.createElement("INPUT");
        x.setAttribute("type", "text");
        x.name = "project_id";
        x.value = myId;

        console.log(x.value);
        e.preventDefault();
    })
})()