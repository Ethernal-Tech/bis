USE [BIS3]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [Owner], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 1, 1, N'549300NROGNBV2T1GS07', N'MY', N'SG', N'100000000', 0, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
