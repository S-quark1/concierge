const navbar = document.querySelector("#accordionSidebarMain");
const content = document.querySelector("#content");

navbar.addEventListener("click", e => {
    const selected = e.target.id;
    const xhr = new XMLHttpRequest();
    xhr.open("GET", `/my-cabinet-admin/${selected}`, true);
    console.log(selected)
    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            content.innerHTML = xhr.responseText;
        }
    };
    xhr.send();
});
