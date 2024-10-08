//
//  FloatingButton.swift
//  Zettelgarden
//
//  Created by Nicholas Savage on 2024-08-18.
//

import SwiftUI

struct FloatingButton: View {
    var action: () -> Void
    var imageText: String
    var body: some View {

        Button(
            action: {
                action()
            },
            label: {
                Image(systemName: imageText)
                    .font(.title.weight(.semibold))
                    .padding()
                    .background(Color.blue)
                    .foregroundColor(.white)
                    .clipShape(Circle())
                    .shadow(radius: 4, x: 0, y: 4)
            }
        )
        .padding(7)
    }
}
struct FloatingButton_Preview: PreviewProvider {
    static var previews: some View {
        FloatingButton(
            action: {
                // Action to be triggered when the button is tapped
                print("Floating Button Pressed")
            },
            imageText: "plus"
        )
        .previewLayout(.sizeThatFits)
        .padding()  // Add padding for better appearance in the preview
    }
}
