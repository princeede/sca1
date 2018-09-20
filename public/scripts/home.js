(function (){
    var myId;
    buttonTestArr = $(".supportButton").click(function(event){
        myId = event.currentTarget.dataset["some"]
        console.log(myId)
    });

    var myForm =  document.querySelector(".supportForm");
    console.log(myForm)
    myForm.addEventListener("submit", e => {
        var x = document.createElement("INPUT");
        x.setAttribute("type", "number");
        x.name = "project_id";
        x.value = myId;
        e.currentTarget.insertBefore(x, e.currentTarget.childNodes[0])
        console.log(e.currentTarget);
        // e.preventDefault();
    })
})()