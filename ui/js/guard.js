var li_items = document.querySelectorAll(".sidebar ul li");
var hamburger = document.querySelector(".hamburger");
var wrapper = document.querySelector(".wrapper");

li_items.forEach((li_item)=>{
	li_item.addEventListener("mouseenter", ()=>{

			li_item.closest(".wrapper").classList.remove("hover_collapse");

	})
 
})

li_items.forEach((li_item)=>{
	li_item.addEventListener("mouseleave", ()=>{

			li_item.closest(".wrapper").classList.add("hover_collapse");
ed

	})
})

document.addEventListener('DOMContentLoaded', function() {
    const hamburger = document.querySelector('.hamburger');
    const wrapper = document.querySelector('.wrapper');

    hamburger.addEventListener('click', function() {
      wrapper.classList.toggle('hover_collapse');
      wrapper.classList.toggle('click_collapse');
    });
  });
