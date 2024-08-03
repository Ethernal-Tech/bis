USE [BIS1]
GO
SET IDENTITY_INSERT [dbo].[BankType] ON

INSERT [dbo].[BankType]
    ([Id], [Name])
VALUES
    (1, N'Commercial')
INSERT [dbo].[BankType]
    ([Id], [Name])
VALUES
    (2, N'Central')
SET IDENTITY_INSERT [dbo].[BankType] OFF
GO
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'AU', N'Australia')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'KR', N'South Korea')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'MY', N'Malaysia')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'SG', N'Singapore')
GO
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'54930035WQZLGC45RZ35', N'Monetary Authority of Singapore', N'10 Shenton Wy', N'SG', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'549300NROGNBV2T1GS07', N'Bank Negara Malaysia', N'Jalan Dato Onn', N'MY', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'XXXXXXXXXXXXXXXXXXXX', N'MY Commercial Bank', N'MY Address 1', N'MY', 1)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'YYYYYYYYYYYYYYYYYYYY', N'SG Commercial Bank', N'SG Address 1', N'SG', 1)
GO
SET IDENTITY_INSERT [dbo].[BankClient] ON

INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (1, N'111', N'Company A', N'', N'YYYYYYYYYYYYYYYYYYYY')
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (2, N'222', N'Company B', N'', N'XXXXXXXXXXXXXXXXXXXX')
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (3, N'333', N'Company C', N'', N'YYYYYYYYYYYYYYYYYYYY')
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (4, N'444', N'Company D', N'', N'XXXXXXXXXXXXXXXXXXXX')
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (5, N'555', N'Company E', N'', N'YYYYYYYYYYYYYYYYYYYY')
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (6, N'666', N'Company F', N'', N'XXXXXXXXXXXXXXXXXXXX')
SET IDENTITY_INSERT [dbo].[BankClient] OFF
GO
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'1', N'SG Commercial Bank Admin', N'SGCadmin', N'admin', N'YYYYYYYYYYYYYYYYYYYY')
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'2', N'Monetary Authority of Singapore Admin', N'MASadmin', N'admin', N'54930035WQZLGC45RZ35')
GO
SET IDENTITY_INSERT [dbo].[PolicyType] ON

INSERT [dbo].[PolicyType]
    ([Id], [Code], [Name])
VALUES
    (1, N'CFM', N'Capital Flow Management')
INSERT [dbo].[PolicyType]
    ([Id], [Code], [Name])
VALUES
    (2, N'SCL', N'Sanctions')
SET IDENTITY_INSERT [dbo].[PolicyType] OFF
GO
SET IDENTITY_INSERT [dbo].[TransactionType] ON

INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (1, N'DDWN', N'Loan Drawdown')
INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (2, N'PPAY', N'Loan Repayment')
SET IDENTITY_INSERT [dbo].[TransactionType] OFF
GO
SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [Owner], [TransactionTypeId], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [BeneficiaryJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 1, N'549300NROGNBV2T1GS07', 1, N'MY', N'SG', N'MY', N'100000000', 0, 1)
SET IDENTITY_INSERT [dbo].[Policy] OFF
GO
INSERT [dbo].[Transaction]
    ([Id], [OriginatorBankId], [BeneficiaryBankId], [SenderId], [ReceiverId], [Currency], [Amount], [TransactionTypeId], [LoanId])
VALUES
    (N'6fa5d15ac6d412a9a0bcb64c3c42dfa362d0dd2fa8e80f5f8a5aa2233256e458', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 3, 4, N'SGD', 200000, 1, 7920367)
INSERT [dbo].[Transaction]
    ([Id], [OriginatorBankId], [BeneficiaryBankId], [SenderId], [ReceiverId], [Currency], [Amount], [TransactionTypeId], [LoanId])
VALUES
    (N'96a490b0b49084a2b2d0358ad468ef39601fda394e1c28b8cb6c2a5573c2c362', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 1, 2, N'SGD', 100000, 1, 8199270)
INSERT [dbo].[Transaction]
    ([Id], [OriginatorBankId], [BeneficiaryBankId], [SenderId], [ReceiverId], [Currency], [Amount], [TransactionTypeId], [LoanId])
VALUES
    (N'a98b630b1147b832beeb6af36922e358f4dddd03c587fbb59a24eecc92a4645a', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 5, 6, N'SGD', 300000, 1, 5748927)
INSERT [dbo].[Transaction]
    ([Id], [OriginatorBankId], [BeneficiaryBankId], [SenderId], [ReceiverId], [Currency], [Amount], [TransactionTypeId], [LoanId])
VALUES
    (N'fb2f320a369c9a73e15b48d1a08c33dd0a6876ebd24a357a32d9ea90972fca06', N'XXXXXXXXXXXXXXXXXXXX', N'YYYYYYYYYYYYYYYYYYYY', 2, 1, N'SGD', 50000, 2, 2789161)
GO
INSERT [dbo].[TransactionPolicy]
    ([TransactionId], [PolicyId], [Status], [AdditionalParameters], [Description])
VALUES
    (N'6fa5d15ac6d412a9a0bcb64c3c42dfa362d0dd2fa8e80f5f8a5aa2233256e458', 1, 1, N'', N'')
INSERT [dbo].[TransactionPolicy]
    ([TransactionId], [PolicyId], [Status], [AdditionalParameters], [Description])
VALUES
    (N'96a490b0b49084a2b2d0358ad468ef39601fda394e1c28b8cb6c2a5573c2c362', 1, 1, N'', N'')
INSERT [dbo].[TransactionPolicy]
    ([TransactionId], [PolicyId], [Status], [AdditionalParameters], [Description])
VALUES
    (N'a98b630b1147b832beeb6af36922e358f4dddd03c587fbb59a24eecc92a4645a', 1, 1, N'', N'')
GO
SET IDENTITY_INSERT [dbo].[Status] ON

INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (1, N'TransactionCreated')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (2, N'PoliciesApplied')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (3, N'ComplianceProofRequested')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (4, N'ComplianceCheckPassed')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (5, N'ProofInvalid')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (6, N'AssetSent')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (7, N'TransactionCompleted')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (8, N'TransactionCanceled')
SET IDENTITY_INSERT [dbo].[Status] OFF
GO
INSERT [dbo].[TransactionHistory]
    ([TransactionId], [StatusId], [Date])
VALUES
    (N'6fa5d15ac6d412a9a0bcb64c3c42dfa362d0dd2fa8e80f5f8a5aa2233256e458', 1, CAST(N'2024-07-24T12:31:28.5536046' AS DateTime2))
INSERT [dbo].[TransactionHistory]
    ([TransactionId], [StatusId], [Date])
VALUES
    (N'96a490b0b49084a2b2d0358ad468ef39601fda394e1c28b8cb6c2a5573c2c362', 1, CAST(N'2024-07-24T12:29:53.4445250' AS DateTime2))
INSERT [dbo].[TransactionHistory]
    ([TransactionId], [StatusId], [Date])
VALUES
    (N'a98b630b1147b832beeb6af36922e358f4dddd03c587fbb59a24eecc92a4645a', 1, CAST(N'2024-07-24T12:33:01.3220043' AS DateTime2))
INSERT [dbo].[TransactionHistory]
    ([TransactionId], [StatusId], [Date])
VALUES
    (N'fb2f320a369c9a73e15b48d1a08c33dd0a6876ebd24a357a32d9ea90972fca06', 1, CAST(N'2024-07-24T12:35:33.8254274' AS DateTime2))
GO
SET IDENTITY_INSERT [dbo].[TransactionProof] ON

INSERT [dbo].[TransactionProof]
    ([Id], [TransactionId], [Proof])
VALUES
    (1, N'96a490b0b49084a2b2d0358ad468ef39601fda394e1c28b8cb6c2a5573c2c362', N'96a490b0b49084a2b2d0358ad468ef39601fda394e1c28b8cb6c2a5573c2c362,0;0xaD472ffC8a7a59db1dd4367cbb2B21A1c1A798F7;90c8d606fd4c1ae69c78e0fc9ba4dd7e96d40af6d15448fc1677bf2fef10edfc5eef534b6ab0168eed940e2bfe44971134f8c0208e4aed29fdfbed7b4feab0da1b')
INSERT [dbo].[TransactionProof]
    ([Id], [TransactionId], [Proof])
VALUES
    (2, N'6fa5d15ac6d412a9a0bcb64c3c42dfa362d0dd2fa8e80f5f8a5aa2233256e458', N'6fa5d15ac6d412a9a0bcb64c3c42dfa362d0dd2fa8e80f5f8a5aa2233256e458,0;0x918ED453521250764F4a5db1af8213d3Bd10Ef64;58b43c01692b55f7a537651322c8e7f7f2eb9c1a3eaa59f9dcd4793882f3b5872282178d2ca7792e79750b471cd4edd47770f4bc1d38f854a615cf3ec1b62aad1b')
INSERT [dbo].[TransactionProof]
    ([Id], [TransactionId], [Proof])
VALUES
    (3, N'a98b630b1147b832beeb6af36922e358f4dddd03c587fbb59a24eecc92a4645a', N'a98b630b1147b832beeb6af36922e358f4dddd03c587fbb59a24eecc92a4645a,0;0x413893fd62F0058803a47aF421e00d0B4Ec31b89;b33ebe53ad02943821a164f0d7305e9ea2da9439ca6ed3dd4c0c5aac021c21017e9999e14fc4e654097424adab28b0f17744e88442046640076ac117e5b982681b')

-- New Transactions
INSERT [dbo].[Transaction]
    ([Id], [OriginatorBankId], [BeneficiaryBankId], [SenderId], [ReceiverId], [Currency], [Amount], [TransactionTypeId], [LoanId])
VALUES
    (N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 1, 3, N'SGD', 250000, 1, 1234567),
    (N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd', N'XXXXXXXXXXXXXXXXXXXX', N'YYYYYYYYYYYYYYYYYYYY', 4, 5, N'SGD', 150000, 2, 2345678),
    (N'1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 1, 2, N'SGD', 210000, 1, 3456789),
    (N'2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f', N'XXXXXXXXXXXXXXXXXXXX', N'YYYYYYYYYYYYYYYYYYYY', 3, 4, N'SGD', 220000, 2, 4567890),
    (N'3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 5, 6, N'SGD', 230000, 1, 5678901),
    (N'4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h', N'XXXXXXXXXXXXXXXXXXXX', N'YYYYYYYYYYYYYYYYYYYY', 2, 1, N'SGD', 240000, 2, 6789012),
    (N'5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 1, 3, N'SGD', 250000, 1, 7890123),
    (N'6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j', N'XXXXXXXXXXXXXXXXXXXX', N'YYYYYYYYYYYYYYYYYYYY', 4, 5, N'SGD', 260000, 2, 8901234),
    (N'7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 5, 6, N'SGD', 270000, 1, 9012345),
    (N'8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l', N'XXXXXXXXXXXXXXXXXXXX', N'YYYYYYYYYYYYYYYYYYYY', 2, 1, N'SGD', 280000, 2, 1234567),
    (N'9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m', N'YYYYYYYYYYYYYYYYYYYY', N'XXXXXXXXXXXXXXXXXXXX', 1, 3, N'SGD', 290000, 1, 2345678),
    (N'0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m0n', N'XXXXXXXXXXXXXXXXXXXX', N'YYYYYYYYYYYYYYYYYYYY', 4, 5, N'SGD', 300000, 2, 3456789);
-- Transaction Policies for New Transactions
INSERT [dbo].[TransactionPolicy]
    ([TransactionId], [PolicyId], [Status], [AdditionalParameters], [Description])
VALUES
    (N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd', 1, 1, N'', N''),
    (N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd', 1, 1, N'', N''),
    (N'1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e', 1, 1, N'', N''),
    (N'2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f', 1, 1, N'', N''),
    (N'3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g', 1, 1, N'', N''),
    (N'4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h', 1, 1, N'', N''),
    (N'5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i', 1, 1, N'', N''),
    (N'6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j', 1, 1, N'', N''),
    (N'7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k', 1, 1, N'', N''),
    (N'8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l', 1, 1, N'', N''),
    (N'9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m', 1, 1, N'', N''),
    (N'0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m0n', 1, 1, N'', N'');

-- Transaction Histories for New Transactions
INSERT [dbo].[TransactionHistory]
    ([TransactionId], [StatusId], [Date])
VALUES
    (N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd', 1, GETDATE()),
    (N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd', 1, GETDATE()),
    (N'1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e', 1, GETDATE()),
    (N'2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f', 1, GETDATE()),
    (N'3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g', 1, GETDATE()),
    (N'4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h', 1, GETDATE()),
    (N'5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i', 1, GETDATE()),
    (N'6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j', 1, GETDATE()),
    (N'7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k', 1, GETDATE()),
    (N'8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l', 1, GETDATE()),
    (N'9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m', 1, GETDATE()),
    (N'0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m0n', 1, GETDATE());

-- Transaction Proofs for New Transactions
INSERT [dbo].[TransactionProof]
    ([Id], [TransactionId], [Proof])
VALUES
    (4, N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd', N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd,0;0x1234567890abcdef1234567890abcdef12345678;abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd'),
    (5, N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd', N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd,0;0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef;1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd'),
    (6, N'1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e', N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd,0;0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd'),
    (7, N'2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f', N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd,0;0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd'),
    (8, N'3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g', N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd,0;0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd'),
    (9, N'4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h', N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd,0;0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd'),
    (10, N'5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i', N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd,0;0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd'),
    (11, N'6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j', N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd,0;0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd'),
    (12, N'7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k', N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd,0;0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd'),
    (13, N'8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l', N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd,0;0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd'),
    (14, N'9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m', N'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd,0;0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd'),
    (15, N'0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9m0n', N'1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd,0;0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefabcd');


SET IDENTITY_INSERT [dbo].[TransactionProof] OFF
GO

