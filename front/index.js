let tg = window.Telegram.WebApp;

window.addEventListener("DOMContentLoaded", (event) => {
    let btn = document.getElementById("my_kek_button");

    btn.addEventListener("click", function() {
        var captain_name = document.getElementById("captain_name")
        var group = document.getElementById("group")
        var phone = document.getElementById("phone")
        var team_name = document.getElementById("team_name")
        var team_size = document.getElementById("team_size")
        var captain_name = document.getElementById("captain_name")

        fetch("http://localhost:8000/register", {
            method: "POST",
            body: JSON.stringify({
                userId: 1,
                captain_name: captain_name,
                group: group,
                phone: phone,
                team_name: team_name,
                team_size: team_size,
                cpatain_name: captain_name,
            }),
            headers: {
                "Content-type": "application/json; charset=UTF-8"
            }
        })
            .then((response) => response.json())
            .then((json) => console.log(json));
    });
});
