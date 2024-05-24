USE BIS
GO

ALTER TABLE [dbo].[Bank] ADD [BankTypeId] [int]

ALTER TABLE [dbo].[Bank]  WITH CHECK ADD  CONSTRAINT [FK_Bank_BankType] FOREIGN KEY([BankTypeId])
REFERENCES [dbo].[BankType] ([Id])
GO

ALTER TABLE [dbo].[Bank] CHECK CONSTRAINT [FK_Bank_BankType]
GO

UPDATE [dbo].[Bank]
SET BankTypeId = 1
GO

ALTER TABLE [dbo].[Bank]
ALTER COLUMN [BankTypeId] [int] NOT NULL;
GO