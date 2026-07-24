package db

import (
	"sync"
	"time"

	"finedu-backend/internal/models"
)

var (
	FlashcardStore   []models.Flashcard
	FlashcardSeq     int64 = 1
	FlashcardStoreMu sync.Mutex
)

func InitMockFlashcards() {
	FlashcardStoreMu.Lock()
	defer FlashcardStoreMu.Unlock()

	now := time.Now().UTC()

	FlashcardStore = []models.Flashcard{
		{
			ID:                  "fc-1",
			Title:               "What is compound interest?",
			Category:            "Investing",
			WhyItMatters:        "Compound interest is the foundation of long-term wealth building. Even small amounts invested early can grow exponentially over decades.",
			Definition:          "Compound interest is interest earned on both the principal and previously earned interest, creating exponential growth.",
			Example:             "Invest $1,000 at 8% annual return: Year 1 = $1,080, Year 10 = $2,158, Year 30 = $10,063.",
			CommonMisconception: "Compound interest only applies to investments. It takes decades to see meaningful returns.",
			ReviewCount:         2,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  "fc-2",
			Title:               "Emergency fund basics",
			Category:            "Budgeting",
			WhyItMatters:        "An emergency fund prevents you from going into debt when unexpected expenses arise, providing financial security and peace of mind.",
			Definition:          "An emergency fund is money set aside specifically for unexpected financial emergencies, typically 3-6 months of living expenses.",
			Example:             "If your monthly expenses are $3,000, aim for an emergency fund of $9,000-$18,000.",
			CommonMisconception: "You need to save your emergency fund before investing. You should prioritize both simultaneously.",
			ReviewCount:         1,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  "fc-3",
			Title:               "Credit score factors",
			Category:            "Credit & Debt",
			WhyItMatters:        "Your credit score affects loan rates, insurance premiums, and even employment opportunities. Maintaining a good score saves thousands in interest.",
			Definition:          "A credit score is a numerical representation of your creditworthiness based on your credit history.",
			Example:             "Score ranges: 300-579 (Poor), 580-669 (Fair), 670-739 (Good), 740-799 (Very Good), 800+ (Excellent).",
			CommonMisconception: "Checking your credit score hurts it. Checking your own score is a soft inquiry and has no impact.",
			ReviewCount:         3,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  "fc-4",
			Title:               "Bull and bear markets",
			Category:            "Markets",
			WhyItMatters:        "Understanding market cycles helps you stay invested during downturns and make better long-term decisions.",
			Definition:          "A bull market is when prices are rising and investor confidence is high. A bear market is when prices fall 20%+ and sentiment is negative.",
			Example:             "The 2008 financial crisis was a severe bear market that lasted ~17 months.",
			CommonMisconception: "You should sell everything during a bear market. History shows the best returns often come after bear markets end.",
			ReviewCount:         0,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  "fc-5",
			Title:               "Retirement account types",
			Category:            "Retirement",
			WhyItMatters:        "Choosing the right retirement account can save you tens of thousands in taxes over your lifetime.",
			Definition:          "Retirement accounts are tax-advantaged investment accounts designed to help you save for retirement (401k, IRA, Roth IRA, etc.).",
			Example:             "Traditional IRA contributions reduce taxable income now. Roth IRA contributions grow tax-free.",
			CommonMisconception: "You can only contribute to one retirement account. You can contribute to multiple accounts in the same year.",
			ReviewCount:         2,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  "fc-6",
			Title:               "Tax-loss harvesting",
			Category:            "Taxes",
			WhyItMatters:        "Tax-loss harvesting can reduce your tax liability by thousands each year, improving net returns on your portfolio.",
			Definition:          "Tax-loss harvesting is selling losing investments to offset capital gains and reduce taxable income.",
			Example:             "Sell $5,000 loss to offset $5,000 gain, reducing capital gains tax from $750 (15% rate) to $0.",
			CommonMisconception: "Tax-loss harvesting is complicated and only for rich people. It's a simple strategy anyone can use.",
			ReviewCount:         1,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
	}

	FlashcardSeq = int64(len(FlashcardStore) + 1)
}
