//
//  Dates.swift
//  Zettelgarden
//
//  Created by Nicholas Savage on 2024-07-10.
//

import Foundation

func formatDate(input: Date) -> String {
    if isToday(maybeDate: input) {
        return "Today"
    }
    let dateFormatter = DateFormatter()
    dateFormatter.dateFormat = "YYYY-MM-dd"
    let result = dateFormatter.string(from: input)
    return result
}

func isToday(maybeDate: Date?) -> Bool {
    guard let date = maybeDate else {
        return false
    }

    let calendar = Calendar.current
    let today = Date()

    let todayComponents = calendar.dateComponents([.year, .month, .day], from: today)
    let dateComponents = calendar.dateComponents([.year, .month, .day], from: date)

    return todayComponents.year == dateComponents.year
        && todayComponents.month == dateComponents.month
        && todayComponents.day == dateComponents.day
}