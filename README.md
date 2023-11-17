<div align='center'>
  <h1>Intermediary Web Server</h1>
  <p>
    Intermediary (Golang) Web Server represents a central component in the security checks carried out between banks. Checks are carried out by setting policies such as the total amount of funds that can be sent and checking whether the entity is on the sanction list.
  </p>
</div>

</br>
</br>

<div align='center'>
  <img src='https://ethernal.tech/static/media/ethernal.e8296ae3d14edef13517cc8beed9ad35.svg' width='10%' />
</div>

</br>

<div align='center'>
  Powered by <a href='https://ethernal.tech/'>Ethernal-Tech</a>
</div>

</br>
</br>
</br>

### 🏗 — Architecture (Application structure)
___
                                           
A classic web server that, in addition to its basic functionalities, communicates with the [GPJC-API](https://github.com/Ethernal-Tech/gpjc-api) to verify the list of sanctioned entities. The following techology stack were used in the implementation of the solution:
1. Golang (GO)
2. HTML
3. CSS
4. JavaScript
5. SQL (Microsoft SQL Server)

<br/>

### ⚡ — Requirements
___

There are three requirements to run the application. They are as follows:
1. Golang (GO) compiler - 1.19+ version (<a href='https://go.dev/doc/install'>go download</a>)
2. Microsoft SQL server (<a href='https://learn.microsoft.com/en-us/sql/database-engine/install-windows/install-sql-server?view=sql-server-ver16'>Installation guide</a>)
3. <a href='https://github.com/Ethernal-Tech/gpjc-api'>GPJC-API</a> (follow the instructions in README.md)

<br/>

### ⚙ — Installing and Running applications
___

(for windows OS)

Installation is very simple and consists of only a few steps:
1. Download (clone) the project
2. Create a new Database with the name  `BIS` (**use Windows Authentication**)
3. From the project take the scripts (`../DB/scripts/`) and run them in the following order:
   - `CreateBISdb.sql`
   - `InsertData.sql`
4. Run the previously mentioned GPJC-API
5. Modify the `.env` file (in case of running on multiple machines)
6. Run the `go mod tidy` Golang command
7. Run the Golang Server (`go run .`)

<br/>
<br/>

<p align="center"><a href="https://github.com/Ethernal-Tech/bis#"><img src="http://randojs.com/images/backToTopButtonTransparentBackground.png" alt="Back to top" height="29"/></a></p>
