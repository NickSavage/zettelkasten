//
//  Dates.swift
//  Zettelgarden
//
//  Created by Nicholas Savage on 2024-07-10.
//

import Foundation

func isToday(maybeDate: Date?) -> Bool {
    if let date = maybeDate {
        let calendar = Calendar.current
        return calendar.isDateInToday(date)
    }
    else {
        return false
    }
}
