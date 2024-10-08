//
//  Card.swift
//  Zettelgarden
//
//  Created by Nicholas Savage on 2024-05-13.
//
import Combine
import Foundation
import SwiftUI

public struct Card: Identifiable, Codable {
    public var id: Int
    public var card_id: String
    public var user_id: Int
    public var title: String
    public var body: String
    public var link: String
    public var created_at: Date
    public var updated_at: Date
    public var parent_id: Int
    public var parent: PartialCard?
    //var card_links: [PartialCard]
    public var children: [PartialCard]
    public var references: [PartialCard]
    public var files: [File]
    public var tags: [Tag]
    public var is_flashcard: Bool

    public enum CodingKeys: String, CodingKey {
        case id
        case card_id
        case user_id
        case title
        case body
        case link
        case created_at
        case updated_at
        case parent_id
        case parent
        case children
        case references
        case files
        case tags
        case is_flashcard
    }
    public init(from decoder: Decoder) throws {
        let container: KeyedDecodingContainer<Card.CodingKeys> = try decoder.container(
            keyedBy: CodingKeys.self
        )
        id = try container.decode(Int.self, forKey: .id)
        card_id = try container.decode(String.self, forKey: .card_id)
        user_id = try container.decode(Int.self, forKey: .user_id)
        title = try container.decode(String.self, forKey: .title)
        body = try container.decode(String.self, forKey: .body)
        link = try container.decode(String.self, forKey: .link)
        let createdAtString = try container.decode(
            String.self,
            forKey: .created_at
        )
        created_at = parseDate(input: createdAtString) ?? Date()
        let updatedAtString = try container.decode(
            String.self,
            forKey: .updated_at
        )
        updated_at = parseDate(input: updatedAtString) ?? Date()
        parent_id = try container.decode(Int.self, forKey: .parent_id)
        parent = try container.decodeIfPresent(PartialCard.self, forKey: .parent)
        children = try container.decodeIfPresent([PartialCard].self, forKey: .children) ?? []
        references = try container.decodeIfPresent([PartialCard].self, forKey: .references) ?? []
        files = try container.decodeIfPresent([File].self, forKey: .files) ?? []
        tags = try container.decodeIfPresent([Tag].self, forKey: .tags) ?? []
        is_flashcard = try container.decodeIfPresent(Bool.self, forKey: .is_flashcard) ?? false
    }
    public init(
        id: Int,
        card_id: String,
        user_id: Int,
        title: String,
        body: String,
        link: String,
        created_at: Date,
        updated_at: Date,
        parent_id: Int,
        parent: PartialCard?,
        children: [PartialCard],
        references: [PartialCard],
        files: [File],
        tags: [Tag],
        is_flashcard: Bool
    ) {
        self.id = id
        self.card_id = card_id
        self.user_id = user_id
        self.title = title
        self.body = body
        self.link = link
        self.created_at = created_at
        self.updated_at = updated_at
        self.parent_id = parent_id
        self.parent = parent
        self.children = children
        self.references = references
        self.files = files
        self.tags = tags
        self.is_flashcard = is_flashcard
    }
}

extension Card {
    public static var sampleData: [Card] =
        [
            Card(
                id: 0,
                card_id: "1",
                user_id: 1,
                title: "hello world",
                body: "this is a test of the emergency response system",
                link: "",
                created_at: Date(),
                updated_at: Date(),
                parent_id: 0,
                parent: nil,
                children: [],
                references: [],
                files: [],
                tags: [],
                is_flashcard: false
            ),
            Card(
                id: 1,
                card_id: "1/A",
                user_id: 1,
                title: "update",
                body: "this is another test of the emergency response system",
                link: "",
                created_at: Date(),
                updated_at: Date(),

                parent_id: 1,
                parent: nil,
                children: [],
                references: [],
                files: [],
                tags: [],
                is_flashcard: false
            ),
        ]

    public static var emptyCard: Card {
        Card(
            id: -1,
            card_id: "",
            user_id: -1,
            title: "",
            body: "",
            link: "",
            created_at: Date(),
            updated_at: Date(),
            parent_id: -1,
            parent: nil,
            children: [],
            references: [],
            files: [],
            tags: [],
            is_flashcard: false
        )
    }
}
