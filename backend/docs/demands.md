## Mess registration docs

Last Updated: Sep 25,2025

### requirements

we have talked about students and mess staff side and you all know it well.

Hostel Office (should remain in power with all logs)

1. swapping
  - they need history for any student swaps (for billing reasons)
  - limit swapping to only once per month for all students
  - keep track of swapping date+time
2. Only one time scanning per meal
  - if scanned again mess people can see last scanned X minutes ago (in yellow color)
3. Hostel office
  - query student by rollno and get his all info
  - can do anything to a student -> deactivate his mess, swap mess.
    - such logs are also to be kept
  -

how it works:

hostel office gets a list of students who can register (around mid of month)
now only those are allowed to register

then first veg registrations start

then later mess registration for all is started

this is done prior to upcoming month as you all know

then from next month onwards they shift current regirstations

EXTRA

4. Mess Rebate

5. Mess Safety

---


## Tables
1. users
2. mess_registrations
3. swap_requests
4. meal_scans
5. monthly_swap_tracking
6. Create a table for loggs.

## Api list

- Auth
  - POST  /api/auth/login
  - POST /api/auth/logout

- Student API's
  - GET  /api/students/profile
  - PUT /api/students/profile
  - POST /api/students/register
  - SWAPS
    - GET /api/students/swaps
    - POST /api/students/swaps
    - POST /api/students/swaps/accept
    - DELETE /api/students/swaps

- Mess Staff
  - GET /api/staff/scan

- Hostel Office API's
  - GET /api/office/students
  - GET /api/office/students/:id
  - PUT /api/office/students/:id - Update any field , can also deactivate or activate
  - GET /api/office/registrations
