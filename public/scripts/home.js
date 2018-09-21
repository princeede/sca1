(function (){
    var myId;
    var myAction;
    buttonTestArr = $(".supportButton").click(function(event){
        myId = event.currentTarget.dataset["some"]
        myAction = event.currentTarget.dataset["action"]
    });

    var myForm =  document.querySelector(".supportForm");
    myForm.addEventListener("submit", e => {
        var id = document.createElement("INPUT");
        id.setAttribute("type", "number");
        id.name = "project_id";
        id.value = myId;
        e.currentTarget.insertBefore(id, e.currentTarget.childNodes[3])

        var action = document.createElement("input")
        action.setAttribute("text", myAction)
        action.name = "action";
        action.value = myAction;
        e.currentTarget.insertBefore(action, e.currentTarget.childNodes[3])


    })

})()