// profile menu at dashboard
const profileButton = document.querySelector("#user-menu-button");
const profileMenu = document.querySelector("#profile-menu");

profileButton.addEventListener("click", function () {
  profileMenu.classList.toggle("hidden");
});
