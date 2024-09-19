//
//  ContentView.swift
//  Zettelgarden
//
//  Created by Nicholas Savage on 2024-05-13.
//

import Combine
import SwiftUI
import ZettelgardenShared

struct ContentView: View {
    @State var isMenuOpen: Bool = false
    @Environment(\.scenePhase) private var scenePhase
    @StateObject var cardViewModel = CardViewModel()
    @StateObject var searchViewModel = SearchViewModel()
    @StateObject var partialViewModel = PartialCardViewModel()
    @StateObject var navigationViewModel: NavigationViewModel
    @StateObject var taskListViewModel = TaskListViewModel()
    @StateObject var tagViewModel = TagViewModel()

    init() {
        let cardViewModel = CardViewModel()
        _cardViewModel = StateObject(wrappedValue: cardViewModel)
        _navigationViewModel = StateObject(
            wrappedValue: NavigationViewModel(cardViewModel: cardViewModel)
        )
    }

    var body: some View {
        NavigationView {
            VStack {

                if navigationViewModel.selection == .tasks {
                    TaskListView(taskListViewModel: taskListViewModel)
                }
                else if navigationViewModel.selection == .home {
                    HomeView(
                        cardViewModel: cardViewModel,
                        navigationViewModel: navigationViewModel,
                        partialViewModel: partialViewModel
                    )
                }
                else if navigationViewModel.selection == .card {
                    CardDisplayView(
                        cardListViewModel: partialViewModel,
                        cardViewModel: cardViewModel,
                        navigationViewModel: navigationViewModel
                    )
                }
                else if navigationViewModel.selection == .files {
                    FileListView()
                }
                else if navigationViewModel.selection == .settings {
                    SettingsView()
                }
            }
            .environmentObject(tagViewModel)
            .overlay {
                SidebarView(
                    isMenuOpen: $isMenuOpen,
                    cardViewModel: cardViewModel,
                    navigationViewModel: navigationViewModel,
                    partialViewModel: partialViewModel,
                    taskListViewModel: taskListViewModel
                )
            }
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button(action: {
                        withAnimation {
                            self.isMenuOpen.toggle()
                        }
                    }) {
                        Image(systemName: "sidebar.left")
                    }
                }
            }
            .toolbar {
                ToolbarItemGroup(placement: .bottomBar) {
                    HStack {
                        Button(action: {
                            navigationViewModel.previousVisit()
                        }) {
                            Image(systemName: "chevron.left")
                        }

                        Button(action: {
                            navigationViewModel.nextVisit()
                        }) {
                            Image(systemName: "chevron.right")
                        }
                        Spacer()

                    }

                }
            }
        }
        .onAppear {
            partialViewModel.displayOnlyTopLevel = true
            partialViewModel.loadCards()
            navigationViewModel.visit(page: .tasks)

        }
        .onChange(of: scenePhase) { newPhase in
            partialViewModel.onScenePhaseChanged(to: newPhase)
            taskListViewModel.onScenePhaseChanged(to: newPhase)
        }
    }
}

#Preview {
    ContentView()
}
