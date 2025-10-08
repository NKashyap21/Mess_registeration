# docs for design

current:


/health [GET]
/login [POST]

/students
  /getUser [GET]
  /register [POST]
  /getSwaps [GET]
  /createSwap [POST]
  /deleteSwap [DELETE]
  /acceptSwap [POST]
/messStaff
  /scanning [GET]
/hostelOffice
  /

Tables:
scans
logs - swapping + hostel office

Updates:


need to keep 2 logs table + with date,time ...
  1. for scanning mess cards
  2. user swaps, hostel office changes

1. swapping should be only once per month
2. only one time scanning per meal
  - first time green
  - second time yellow (may be keep some 5 min buf)
  - not this mess red
    - to do this we must classify each scan into
      - breakfast
      - lunch
      - snacks
      - dinner

---
now for new hostel office changes
1. query student info by email
  - also able to change his mess registration status
  - change his mess
  - add students / activate them by email id
  - currently do for iith students
  - etc...
2. upload upcomoing registations sheet
  - also able to add more to it
  - store upcoming registrations in another col in db
3. need to do veg registations before
  - they should not be able to apply on registrations time
4. extend this all to non-iith students / interns etc...
5. hostel office dashboard to see analytics (low priority but easy to make)

---
work:

Saadiq:  in test app. connect to it and make sure everything works(test locally)
Rayaan: test and work on website and make sure everything works
Dhiraj+Kashyap: backend + db work


---
discussions on sep 27 meet
