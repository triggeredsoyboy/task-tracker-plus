// mobile nav at index
const openMobile = document.querySelector("#mobile-button");
const closeMobile = document.querySelector("#close-mobile-button");
const mobileNavbar = document.querySelector("#mobile-navbar");

openMobile.addEventListener("click", function () {
  mobileNavbar.classList.toggle("hidden");
});

closeMobile.addEventListener("click", function () {
  mobileNavbar.classList.toggle("hidden");
});
