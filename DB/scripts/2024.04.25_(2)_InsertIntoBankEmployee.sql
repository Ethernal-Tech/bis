USE [BIS]
GO

INSERT [dbo].[BankEmployee]
    ([Name], [Username], [Password], [BankId])
VALUES
    (N'Bank Negara Malaysia Admin', N'BNMadmin', N'admin', 3)

INSERT [dbo].[BankEmployee]
    ([Name], [Username], [Password], [BankId])
VALUES
    (N'Monetary Authority of Singapore Admin', N'MASadmin', N'admin', 4)

GO