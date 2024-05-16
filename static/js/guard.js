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

  // Autofill date and time fields with current date and time
	window.onload = function() {
		var dateInput = document.getElementById('date');
		var timeInput = document.getElementById('time');
		var currentDate = new Date();
		var year = currentDate.getFullYear();
		var month = ('0' + (currentDate.getMonth() + 1)).slice(-2);
		var day = ('0' + currentDate.getDate()).slice(-2);
		var hours = ('0' + currentDate.getHours()).slice(-2);
		var minutes = ('0' + currentDate.getMinutes()).slice(-2);
		var formattedDate = year + '-' + month + '-' + day;
		var formattedTime = hours + ':' + minutes;
		dateInput.value = formattedDate;
		timeInput.value = formattedTime;
	};