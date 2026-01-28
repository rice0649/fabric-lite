See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.


### **Compartment 1: Core Data & Personnel Management**

*   **Objective:** I will establish the foundational database in MS Access to act as the single source of truth for all personnel and training data, focusing on the UTM's role in personnel oversight.
*   **To-Do List:**
    1.  **Analyze Data Structures:** I will analyze `tbl_100_Alpha`, `tbl_200_Ojt`, and the other table metadata to design the core data tables in MS Access. The design will be optimized for tracking training-specific information relevant to the UTM.
    2.  **Develop In-Processing Module:** I will develop a VBA-driven form in MS Access to streamline the in-processing of new members. This form will capture all data points required by the "UTM In-processing Interview Template" (`UTM Procedures PSDG.pdf`, Attachment 2) and the maintenance-specific in-processing steps (`Mx_Training_PSD_Guide.pdf`).
    3.  **Create Data Import/Scrubbing Process:** I will create a repeatable process using Access and VBA to import and "scrub" (cleanse) data from source tables like `tbl_100_Alpha` and `tbl_200_Ojt`, ensuring data integrity within our new system.

### **Compartment 2: OJT Roster & Journal Automation**

*   **Objective:** I will automate the creation, annotation, and distribution of the monthly On-the-Job Training (OJT) Roster and integrate the new `tbl_JournalEntries` workflow.
*   **To-Do List:**
    1.  **Build OJT Data Query:** I will build a query in MS Access that pulls and formats the raw OJT data as specified in the `MilPDS OJT Roster Raw Data` (`UTM Procedures PSDG.pdf`, Attachment 4).
    2.  **Develop Annotation Interface:** I will develop an Excel template with a VBA user interface that imports the Access data. This interface will allow UTMs to easily add required monthly annotations as described in the `OJT_Roster_Workflow_Context.txt`.
    3.  **Integrate Journal Entries:** I will create a form in Access for entering data into `tbl_JournalEntries`, linking it to the main personnel table via SSAN. This will replace the manual entry process.
    4.  **Automate Reporting Workflow:** I will write a Power Automate flow that triggers before the UTA, automatically generates a PDF of the annotated roster, and emails it to the necessary leadership.

### **Compartment 3: Career Development Course (CDC) Tracking**

*   **Objective:** I will create a dedicated module to manage the entire lifecycle of CDC administration from the UTM's perspective, including enrollment, progress tracking, and failure processing.
*   **To-Do List:**
    1.  **Design CDC Tracking Tables:** I will extend the Access database with tables specifically designed to track CDC enrollments, module completion dates, test scores, and waiver statuses, per the `UTM Career Development Course Administration PSDG.pdf`.
    2.  **Create CDC Management Forms:** I will build Access forms for supervisors and UTMs to log CDC progress reviews, document counseling, and initiate failure packages or waiver requests.
    3.  **Develop Automated Reminders:** I will use Power Automate to monitor CDC enrollment dates and automatically send reminder emails for trainees approaching the 15-month enrollment limit.
    4.  **Template ILPs and IDPs:** I will create standardized Word/Excel templates for the Individual Learning Plan (ILP) and Individual Development Plan (IDP) based on the provided examples, which can be auto-populated from the Access database.

### **Compartment 4: Formal & Maintenance Training Scheduling**

*   **Objective:** I will streamline and centralize the UTM's responsibilities for requesting, tracking, and documenting all formal training, including Craftsman, Maintenance, and Distance Learning courses.
*   **To-Do List:**
    1.  **Create Centralized Tracking Tables:** I will create tables in Access to track all formal training requests, course schedules, Training Line Numbers (TLNs), and completion statuses, consolidating information from the various PSD guides (Formal Training, Craftsman, DL).
    2.  **Build AF Form 898 Module:** I will develop a system using an Access back-end and an Excel/VBA front-end to manage the AF Form 898, Field Training Requirements Scheduling Document process as detailed in the `Mx_Training_PSD_Guide.pdf`.
    3.  **Automate Notification & Routing:** I will use Power Automate to create tasks and notifications for the UTM when training-related emails (like RIPs or cancellations) are received, ensuring timely action.
    4.  **Develop Status Reports:** I will build reports in Access to provide UTMs and leadership with a clear, real-time view of all pending and completed formal training requests.

### **Compartment 5: Reporting & Metrics (TPM) Dashboard**

*   **Objective:** I will automate the data collection and generation of the monthly Training Performance Metrics (TPM) briefing, transforming it from a static report into an interactive dashboard.
*   **To-Do List:**
    1.  **Consolidate TPM Data Points:** I will create a series of complex queries in Access to consolidate all data points required for the TPM briefing, as outlined in the `Training Performance Metrics PSD Guide.pdf` for all component types (Active, ARC, ANG).
    2.  **Connect to Power BI:** I will establish a direct and secure connection from our MS Access database to Power BI.
    3.  **Design Interactive Dashboard:** I will design a Power BI dashboard that not only replicates the sample TPM slide but also allows leadership to interactively filter data by unit, AFSC, TSC, and date ranges.
    4.  **Schedule Automatic Refresh:** I will configure a scheduled data refresh in the Power BI service, ensuring the TPM dashboard is always current with the latest information from the Access database without manual intervention.
